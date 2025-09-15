package health

import (
	"context"
	"log"
	"strings"
	"time"
	domain "users-api/domain/health"
	"users-api/internal/docker"

	"github.com/docker/docker/api/types/container"
)

type Service struct {
	dockerClient *docker.DockerClient
}

func NewService() (*Service, error) {
	log.Println("Iniciando NewService para health check...")
	client, err := docker.NewDockerClient()
	if err != nil {
		log.Printf("Error creando el cliente Docker: %v", err)
		return nil, err
	}
	log.Println("Cliente Docker creado exitosamente")
	return &Service{dockerClient: client}, nil
}

func (s *Service) CheckServices() domain.HealthResponse {
	log.Println("Iniciando CheckServices...")

	if s.dockerClient == nil {
		log.Println("Error: dockerClient es nil")
		return domain.HealthResponse{
			Services:  []domain.ServiceStatus{},
			Timestamp: time.Now().Format(time.RFC3339),
		}
	}

	ctx := context.Background()
	containers, err := s.dockerClient.Client.ContainerList(ctx, container.ListOptions{
		All: true,
	})
	if err != nil {
		return domain.HealthResponse{
			Services:  []domain.ServiceStatus{},
			Timestamp: time.Now().Format(time.RFC3339),
		}
	}
	log.Printf("Found %d containers", len(containers))

	services := []domain.ServiceStatus{}
	// Mapa de nombres de servicios a buscar y sus puertos
	serviceMap := map[string]struct {
		name string
		port string
	}{
		"finalarqsoft2-users-api":   {name: "users-api", port: "8083"},
		"finalarqsoft2-courses-api": {name: "courses-api", port: "8080"},
		"finalarqsoft2-search-api":  {name: "search-api", port: "8082"},
		"inscriptions-api":          {name: "inscriptions-api", port: "8081"},
	}

	for _, container := range containers {
		log.Printf("Checking container: %s (Names: %v)", container.ID[:12], container.Names)

		for containerPrefix, service := range serviceMap {
			// Iterar sobre los nombres del contenedor
			for _, name := range container.Names {
				// Eliminar el "/" inicial del nombre del contenedor
				cleanName := strings.TrimPrefix(name, "/")

				// Verificar si el nombre del contenedor contiene el prefijo buscado
				if strings.Contains(cleanName, containerPrefix) {
					log.Printf("Found service: %s on port %s", service.name, service.port)
					services = append(services, domain.ServiceStatus{
						Name:      service.name,
						Status:    container.State,
						Port:      service.port,
						Container: container.ID[:12],
						Image:     container.Image,
					})
					break
				}
			}
		}
	}

	return domain.HealthResponse{
		Services:  services,
		Timestamp: time.Now().Format(time.RFC3339),
	}

}
