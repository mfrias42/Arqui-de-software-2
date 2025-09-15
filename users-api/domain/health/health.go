package health

type ServiceStatus struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	Port      string `json:"port"`
	Container string `json:"container"` // ID o nombre del contenedor
	Image     string `json:"image"`     // Imagen Docker utilizada
}

type HealthResponse struct {
	Services  []ServiceStatus `json:"services"`
	Timestamp string          `json:"timestamp"`
}
