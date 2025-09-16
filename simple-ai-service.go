package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type AIRequest struct {
	Question string `json:"question"`
	UserID   int    `json:"user_id"`
}

type AIResponse struct {
	Answer    string `json:"answer"`
	Timestamp string `json:"timestamp"`
}

type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
}

type GeminiCandidate struct {
	Content GeminiContent `json:"content"`
}

type ExpenseData struct {
	ID          int     `json:"id"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Date        string  `json:"date"`
	UserID      int     `json:"user_id"`
}

var (
	aiMutex sync.RWMutex
)

func main() {
	r := gin.Default()
	
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "ai-service"})
	})

	api := r.Group("/api/v1")
	aiAPI := api.Group("/ai")
	
	aiAPI.POST("/chat", handleAIChat)

	r.Run(":8086")
}

func handleAIChat(c *gin.Context) {
	var req AIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("JSON binding error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request: %v", err)})
		return
	}

	aiMutex.Lock()
	defer aiMutex.Unlock()

	// Get user's expense data
	expenses, _ := fetchUserExpenses(req.UserID)
	fmt.Printf("Processing AI request for user %d with %d expenses\n", req.UserID, len(expenses))

	// Format expense data for AI
	structuredData := formatExpenseData(expenses)
	
	// Generate AI response
	aiAnswer := generateAIResponse(req.Question, structuredData)
	
	response := AIResponse{
		Answer:    aiAnswer,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}
	
	c.JSON(http.StatusOK, response)
}

func fetchUserExpenses(userID int) ([]ExpenseData, error) {
	// Fetch from expense service
	resp, err := http.Get("http://localhost:8002/api/v1/expenses")
	if err != nil {
		fmt.Printf("HTTP request error: %v\n", err)
		// Return empty slice instead of error to allow AI to work without expense service
		return []ExpenseData{}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("Expense service returned status: %d\n", resp.StatusCode)
		// Return empty slice instead of error
		return []ExpenseData{}, nil
	}

	var result struct {
		Expenses []ExpenseData `json:"expenses"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("JSON decode error: %v\n", err)
		// Return empty slice instead of error
		return []ExpenseData{}, nil
	}

	// Filter by user ID (simplified - in real app, pass user token)
	var userExpenses []ExpenseData
	for _, expense := range result.Expenses {
		if expense.UserID == userID || userID == 1 { // Allow user 1 to see all for demo
			userExpenses = append(userExpenses, expense)
		}
	}

	fmt.Printf("Fetched %d expenses for user %d\n", len(userExpenses), userID)
	return userExpenses, nil
}

func formatExpenseData(expenses []ExpenseData) string {
	if len(expenses) == 0 {
		return "No expense data available."
	}

	// Group by category
	categoryTotals := make(map[string]float64)
	monthlyTotals := make(map[string]float64)
	
	for _, expense := range expenses {
		categoryTotals[expense.Category] += expense.Amount
		
		// Extract month from date
		if len(expense.Date) >= 7 {
			month := expense.Date[:7] // YYYY-MM format
			monthlyTotals[month] += expense.Amount
		}
	}

	var result strings.Builder
	
	// Category breakdown
	result.WriteString("Expense Categories:\n")
	for category, total := range categoryTotals {
		result.WriteString(fmt.Sprintf("- %s: %.2f\n", category, total))
	}
	
	// Monthly breakdown
	result.WriteString("\nMonthly Totals:\n")
	for month, total := range monthlyTotals {
		result.WriteString(fmt.Sprintf("- %s: %.2f\n", month, total))
	}
	
	// Recent transactions
	result.WriteString("\nRecent Transactions:\n")
	recentCount := 5
	if len(expenses) < recentCount {
		recentCount = len(expenses)
	}
	
	for i := 0; i < recentCount; i++ {
		expense := expenses[len(expenses)-1-i] // Get latest first
		result.WriteString(fmt.Sprintf("- %s: %.2f (%s) on %s\n", 
			expense.Description, expense.Amount, expense.Category, expense.Date))
	}

	return result.String()
}

func generateAIResponse(question, data string) string {
	// Try Gemini AI first, fallback to rule-based
	geminiResponse := callGeminiAPI(question, data)
	if geminiResponse != "" {
		return geminiResponse
	}
	
	// Check if it's a general question (not expense-related)
	expenseKeywords := []string{"spend", "spent", "expense", "money", "cost", "budget", "category", "food", "transport", "total", "summary", "financial"}
	isExpenseQuestion := false
	questionLower := strings.ToLower(question)
	
	for _, keyword := range expenseKeywords {
		if strings.Contains(questionLower, keyword) {
			isExpenseQuestion = true
			break
		}
	}
	
	if !isExpenseQuestion {
		// Handle general questions without Gemini
		return handleGeneralQuestion(question)
	}
	
	// Fallback to rule-based responses for expense questions
	return generateRuleBasedResponse(question, data)
}

func handleGeneralQuestion(question string) string {
	question = strings.ToLower(question)
	
	if strings.Contains(question, "hi") || strings.Contains(question, "hello") || strings.Contains(question, "hey") {
		return "Hello! 👋 I'm your FinTrack AI assistant. I can help you analyze your expenses or chat about anything else. What would you like to know? 😊"
	}
	
	if strings.Contains(question, "weather") {
		return "I don't have access to real-time weather data 🌤️, but I can help you track your expenses! You could also ask me about budgeting tips or financial advice 💡"
	}
	
	if strings.Contains(question, "joke") {
		return "Why don't money trees ever grow? Because people keep spending all the seeds! 😄💰 Speaking of money, want to see how you've been spending yours?"
	}
	
	if strings.Contains(question, "ai") || strings.Contains(question, "artificial intelligence") {
		return "AI is fascinating! 🤖 It's technology that can learn and make decisions like humans. I'm an AI assistant built to help you manage your finances better. Want to see what insights I can give about your spending? 📊"
	}
	
	if strings.Contains(question, "how are you") || strings.Contains(question, "how do you do") {
		return "I'm doing great, thanks for asking! 😊 I'm here and ready to help you with your expenses or answer any questions. How can I assist you today? 💡"
	}
	
	if strings.Contains(question, "save money") || strings.Contains(question, "saving") {
		return "Great question! 💸 Here are some money-saving tips:\n• Track all expenses (like you're doing!)\n• Set a monthly budget\n• Cook at home more\n• Compare prices before buying\n• Avoid impulse purchases\nWant me to analyze your current spending patterns? 📊"
	}
	
	if strings.Contains(question, "thank") {
		return "You're very welcome! 😊 I'm always here to help with your finances or any other questions. Feel free to ask me anything! 💡"
	}
	
	// Default response for unknown general questions
	return fmt.Sprintf("That's an interesting question! 🤔 While I specialize in financial management, I'm always happy to chat. I notice you have expense data - would you like me to analyze your spending patterns instead? Or feel free to ask me anything else! 😊")
}

func callGeminiAPI(question, data string) string {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		// No API key, use fallback
		return ""
	}

	// Check if question is expense-related
	expenseKeywords := []string{"spend", "spent", "expense", "money", "cost", "budget", "category", "food", "transport", "total", "summary", "financial"}
	isExpenseQuestion := false
	questionLower := strings.ToLower(question)
	
	for _, keyword := range expenseKeywords {
		if strings.Contains(questionLower, keyword) {
			isExpenseQuestion = true
			break
		}
	}

	var prompt string
	if isExpenseQuestion {
		// Expense-focused prompt
		prompt = fmt.Sprintf(`You are FinTrack, an AI-powered financial assistant.
You help users understand their personal expenses and provide financial insights.

### Your Role:
- Analyze expense data and provide insights
- Give financial advice based on spending patterns
- Help users make better financial decisions
- Be friendly, helpful, and use relevant emojis

### User Question:
%s

### Available Expense Data:
%s

### Instructions:
- If the question is about expenses, use the data provided
- Provide specific numbers and percentages when possible
- Add helpful financial tips and advice
- Use emojis like 💰🍔🚗🛍️📊
- Keep responses clear and actionable

Answer:`, question, data)
	} else {
		// General conversation prompt
		prompt = fmt.Sprintf(`You are FinTrack AI, a friendly and intelligent assistant.
While you specialize in financial management, you can help with various topics.

### Your Personality:
- Friendly, helpful, and conversational
- Smart and knowledgeable about many topics
- Always try to be useful and engaging
- Use appropriate emojis to make conversations fun
- Keep responses concise but informative

### User Question:
%s

### Instructions:
- Answer the question helpfully and accurately
- If it's not finance-related, still be helpful and engaging
- Use a friendly, conversational tone
- Add relevant emojis when appropriate
- If you don't know something, be honest about it

Answer:`, question)
	}

	reqBody := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: prompt},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return ""
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=%s", apiKey)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Gemini API error: %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("Gemini API status: %d\n", resp.StatusCode)
		return ""
	}

	var geminiResp GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		fmt.Printf("Gemini decode error: %v\n", err)
		return ""
	}

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text
	}

	fmt.Println("No valid response from Gemini")
	return ""
}

func generateRuleBasedResponse(question, data string) string {
	question = strings.ToLower(question)
	
	// Parse data for quick access
	lines := strings.Split(data, "\n")
	categories := make(map[string]float64)
	var totalSpent float64
	
	for _, line := range lines {
		if strings.Contains(line, ": ") && strings.HasPrefix(line, "- ") {
			parts := strings.Split(line, ": ")
			if len(parts) == 2 {
				var amount float64
				fmt.Sscanf(parts[1], "%f", &amount)
				category := strings.TrimPrefix(parts[0], "- ")
				categories[category] = amount
				totalSpent += amount
			}
		}
	}

	// Question analysis and responses
	if strings.Contains(question, "total") || strings.Contains(question, "spent") {
		if totalSpent > 0 {
			return fmt.Sprintf("You've spent a total of %.2f 💰. Your biggest expense category is %s.", 
				totalSpent, findLargestCategory(categories))
		}
		return "I don't have any expense information right now 📊."
	}
	
	if strings.Contains(question, "food") || strings.Contains(question, "eating") {
		// Check for both "Food" and "food" categories
		foodAmount := 0.0
		for category, amount := range categories {
			if strings.ToLower(category) == "food" {
				foodAmount += amount
			}
		}
		if foodAmount > 0 {
			return fmt.Sprintf("You spent %.2f on food 🍔. That's %.1f%% of your total expenses.", 
				foodAmount, (foodAmount/totalSpent)*100)
		}
		return "I don't see any food expenses in your data right now 🍽️."
	}
	
	if strings.Contains(question, "transport") || strings.Contains(question, "travel") || strings.Contains(question, "gas") {
		// Check for transport categories (case-insensitive)
		transportAmount := 0.0
		for category, amount := range categories {
			if strings.ToLower(category) == "transport" {
				transportAmount += amount
			}
		}
		if transportAmount > 0 {
			return fmt.Sprintf("You spent %.2f on transport 🚗.", transportAmount)
		}
		return "I don't see any transport expenses in your data right now 🚗."
	}
	
	if strings.Contains(question, "category") || strings.Contains(question, "categories") {
		if len(categories) > 0 {
			result := "Here are your expense categories:\n"
			for category, amount := range categories {
				emoji := getCategoryEmoji(category)
				result += fmt.Sprintf("• %s: %.2f %s\n", category, amount, emoji)
			}
			return result
		}
		return "I don't have any expense categories to show right now 📊."
	}
	
	if strings.Contains(question, "summary") || strings.Contains(question, "overview") {
		if totalSpent > 0 {
			largest := findLargestCategory(categories)
			return fmt.Sprintf("📊 Expense Summary:\n• Total spent: %.2f\n• Categories: %d\n• Largest category: %s (%.2f)\n• You're doing great tracking your expenses! 👍", 
				totalSpent, len(categories), largest, categories[largest])
		}
		return "I don't have enough data for a summary right now 📊."
	}
	
	// Default response
	if totalSpent > 0 {
		return fmt.Sprintf("I can help you analyze your expenses! You've spent %.2f across %d categories. Ask me about specific categories, totals, or summaries 💡.", 
			totalSpent, len(categories))
	}
	
	return "I don't have that information right now. Try asking about your total expenses, food spending, or expense categories! 🤔"
}

func findLargestCategory(categories map[string]float64) string {
	var largest string
	var maxAmount float64
	
	for category, amount := range categories {
		if amount > maxAmount {
			maxAmount = amount
			largest = category
		}
	}
	
	return largest
}

func getCategoryEmoji(category string) string {
	category = strings.ToLower(category)
	
	switch {
	case strings.Contains(category, "food"):
		return "🍔"
	case strings.Contains(category, "transport"):
		return "🚗"
	case strings.Contains(category, "shopping"):
		return "🛍️"
	case strings.Contains(category, "entertainment"):
		return "🎬"
	case strings.Contains(category, "health"):
		return "🏥"
	case strings.Contains(category, "education"):
		return "📚"
	default:
		return "💰"
	}
}