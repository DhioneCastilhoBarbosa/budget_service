package budget

import (
	"budget-service/internal/budget/models"
	"budget-service/internal/budget/s3helper"
	"budget-service/internal/database"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Criar um orçamento (chamado pelo Chat Service)
func CreateBudget(c *gin.Context) {
	// Aumenta o limite do corpo da requisição
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 100<<20)

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao ler formulário multipart: " + err.Error()})
		return
	}

	var budget models.Budget

	// Campos de texto
	budget.SessionID = getFirst(form.Value["session_id"])
	budget.Name = getFirst(form.Value["name"])
	budget.Email = getFirst(form.Value["email"])
	budget.Phone = getFirst(form.Value["phone"])
	budget.LocationType = getFirst(form.Value["location_type"])
	budget.Distance = getFirst(form.Value["distance"])
	budget.NetworkType = getFirst(form.Value["network_type"])
	budget.StructureType = getFirst(form.Value["structure_type"])
	budget.ChargerType = getFirst(form.Value["charge_type"])
	budget.Power = getFirst(form.Value["power"])
	budget.Protection = getFirst(form.Value["protection"])
	budget.Notes = getFirst(form.Value["notes"])
	budget.InstallerName = getFirst(form.Value["installer_name"])

	// Campos numéricos
	if val, err := strconv.ParseUint(getFirst(form.Value["station_count"]), 10, 32); err == nil {
		budget.StationCount = uint(val)
	}
	if val, err := strconv.ParseFloat(getFirst(form.Value["value"]), 64); err == nil {
		budget.Value = val
	}
	installerID := getFirst(form.Value["installer_id"])
	if installerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "installer_id é obrigatório"})
		return
	}
	budget.InstallerID = installerID

	// Upload de fotos com nome seguro
	photoMap := []struct {
		field string
		dest  *string
	}{
		{"photo1", &budget.Photo1},
		{"photo2", &budget.Photo2},
	}

	for _, item := range photoMap {
		files := form.File[item.field]
		if len(files) == 0 {
			continue
		}

		fileHeader := files[0]
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao abrir o arquivo para upload"})
			return
		}

		// Detecta tipo de arquivo
		contentType := fileHeader.Header.Get("Content-Type")
		var ext string
		switch contentType {
		case "image/jpeg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/webp":
			ext = ".webp"
		default:
			ext = ".jpg" // fallback
		}

		// Gera nome seguro
		safeName := fmt.Sprintf("foto_%s_%d%s", item.field, time.Now().UnixNano(), ext)

		// Upload
		url, uploadErr := s3helper.UploadReaderToS3(file, safeName)
		file.Close()

		if uploadErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao enviar %s: %v", item.field, uploadErr)})
			return
		}

		*item.dest = url
	}

	// Salva no banco
	if err := database.DB.Create(&budget).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar orçamento"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Orçamento criado com sucesso",
		"data":    budget,
	})
}

// Auxiliar para obter o primeiro valor de um campo
func getFirst(values []string) string {
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

// Buscar orçamentos do usuário autenticado
func GetUserBudgets(c *gin.Context) {
	// Tenta pegar o user_id via URL param ou query string
	userID := c.Param("user_id")
	if userID == "" {
		userID = c.Query("user_id")
	}

	// Valida se o ID foi fornecido
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetro 'user_id' é obrigatório"})
		return
	}

	var budgets []models.Budget

	// Realiza a busca no banco
	result := database.DB.Where("user_id = ?", userID).Find(&budgets)
	if result.Error != nil {
		log.Printf("Erro ao buscar orçamentos do usuário %s: %v", userID, result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar orçamentos"})
		return
	}

	// Verifica se há resultados
	if len(budgets) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Nenhum orçamento encontrado"})
		return
	}

	// Retorna os orçamentos encontrados
	c.JSON(http.StatusOK, budgets)
}

// busca todos os orçamentos
func GetAllBudgets(c *gin.Context) {
	var budgets []models.Budget

	if err := database.DB.Find(&budgets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar orçamentos"})
		return
	}

	c.JSON(http.StatusOK, budgets)
}

// Vincular orçamentos ao usuário após login
func LinkBudgetsToUser(c *gin.Context) {
	var req struct {
		SessionID string `json:"session_id"`
		UserID    string `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Atualiza apenas orçamentos que ainda não estão vinculados a um usuário
	result := database.DB.Model(&models.Budget{}).
		Where("session_id = ? AND user_id IS NULL", req.SessionID).
		Update("user_id", req.UserID)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusConflict, gin.H{"message": "Orçamentos já vinculados a um usuário"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Orçamentos vinculados ao usuário"})
}

func UpdateBudgetValue(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Value float64 `json:"value"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Valor inválido"})
		return
	}

	if err := database.DB.Model(&models.Budget{}).Where("id = ?", id).Update("value", body.Value).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar valor"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Valor atualizado com sucesso"})
}

// Atualiza o status do orçamento
func UpdateBudgetStatus(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Status string `json:"status"`
	}

	// ✅ Adiciona esta linha!
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	statusList := []string{"aguardando orçamento", "em andamento", "concluido", "cancelado"}

	valid := false
	for _, s := range statusList {
		if s == body.Status {
			valid = true
			break
		}
	}

	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status inválido"})
		return
	}

	if err := database.DB.Model(&models.Budget{}).Where("id = ?", id).Update("status", body.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status atualizado com sucesso"})
}

// Atualiza o data de execução e finalização do orçamento
func UpdateBudgetDates(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		ExecutionDate string `json:"execution_date"`
		FinishDate    string `json:"finish_date"`
	}

	layout := "2006-01-02" // Esperado: "yyyy-mm-dd"
	update := map[string]interface{}{}

	// Valida e converte execution_date
	if body.ExecutionDate != "" {
		t, err := time.Parse(layout, body.ExecutionDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato da execution_date inválido. Use yyyy-mm-dd"})
			return
		}
		update["execution_date"] = t
	}

	// Valida e converte finish_date
	if body.FinishDate != "" {
		t, err := time.Parse(layout, body.FinishDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato da finish_date inválido. Use yyyy-mm-dd"})
			return
		}
		update["finish_date"] = t
	}

	if len(update) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nenhuma data fornecida"})
		return
	}

	if err := database.DB.Model(&models.Budget{}).Where("id = ?", id).Updates(update).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar datas"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Datas atualizadas com sucesso"})
}

// Atualiza o status de pagamento do orçamento
func UpdatePaymentStatus(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		PaymentStatus string `json:"payment_status"`
	}

	if body.PaymentStatus != "pendente" && body.PaymentStatus != "pago" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status de pagamento inválido"})
		return
	}

	if err := database.DB.Model(&models.Budget{}).Where("id = ?", id).Update("payment_status", body.PaymentStatus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar status de pagamento"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status de pagamento atualizado com sucesso"})
}

// Atualiza as confirmações do orçamento
func ConfirmExecution(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		InstallerConfirm *bool `json:"installer_confirm"`
		ClientConfirm    *bool `json:"client_confirm"`
	}

	update := map[string]interface{}{}
	if body.InstallerConfirm != nil {
		update["installer_confirm"] = *body.InstallerConfirm
	}
	if body.ClientConfirm != nil {
		update["client_confirm"] = *body.ClientConfirm
	}

	if len(update) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nenhuma confirmação fornecida"})
		return
	}

	if err := database.DB.Model(&models.Budget{}).Where("id = ?", id).Updates(update).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar confirmações"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Confirmações atualizadas com sucesso"})
}
