package drip

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDrip_Active(t *testing.T) {
	// Configurar un mock de EnvGetter para Drip activo
	mockEnvGetter := &MockEnvGetter{
		ReturnValue: true,
		ReturnFloat: 50.0, // Umbral del 50%
	}

	// Crear un motor de prueba de Gin
	r := gin.Default()

	// Crear un enrutador de prueba y agregar una ruta
	r.POST("/test", func(c *gin.Context) {
		result := Drip(c, mockEnvGetter)
		// Verificar que la función devuelva true (aceptar tráfico)
		assert.True(t, result)
		// Responder con un estado OK
		c.JSON(http.StatusOK, nil)
	})

	// Crear una solicitud HTTP de prueba
	req, _ := http.NewRequest("POST", "/test", nil)
	w := httptest.NewRecorder()

	// Enviar la solicitud al enrutador de prueba
	r.ServeHTTP(w, req)

	// Verificar que el estado de respuesta sea OK (200)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDrip_Inactive(t *testing.T) {
	// Configurar un mock de EnvGetter para Drip inactivo
	mockEnvGetter := &MockEnvGetter{
		ReturnValue: false,
	}

	// Crear un motor de prueba de Gin
	r := gin.Default()

	// Crear un enrutador de prueba y agregar una ruta
	r.POST("/test", func(c *gin.Context) {
		result := Drip(c, mockEnvGetter)
		// Verificar que la función devuelva true (aceptar tráfico)
		assert.True(t, result)
		// Responder con un estado OK
		c.JSON(http.StatusOK, nil)
	})

	// Crear una solicitud HTTP de prueba
	req, _ := http.NewRequest("POST", "/test", nil)
	w := httptest.NewRecorder()

	// Enviar la solicitud al enrutador de prueba
	r.ServeHTTP(w, req)

	// Verificar que el estado de respuesta sea OK (200)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDrip_Drop(t *testing.T) {
	// Configurar un mock de EnvGetter para Drip activo pero por debajo del umbral
	mockEnvGetter := &MockEnvGetter{
		ReturnValue: true,
		ReturnFloat: 10.0, // Umbral del 10%
	}

	// Crear un motor de prueba de Gin
	r := gin.Default()

	// Crear un enrutador de prueba y agregar una ruta
	r.POST("/test", func(c *gin.Context) {
		result := Drip(c, mockEnvGetter)
		// Verificar que la función devuelva false (rechazar tráfico)
		assert.False(t, result)
		// Responder con un estado de servicio no disponible
		c.JSON(http.StatusServiceUnavailable, nil)
	})

	// Crear una solicitud HTTP de prueba
	req, _ := http.NewRequest("POST", "/test", nil)
	w := httptest.NewRecorder()

	// Enviar la solicitud al enrutador de prueba
	r.ServeHTTP(w, req)

	// Verificar que el estado de respuesta sea Service Unavailable (503)
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

type MockEnvGetter struct {
	ReturnValue bool
	ReturnInt   int
	ReturnFloat float64
	ReturnStr   string
}

func (m *MockEnvGetter) GetBoolEnv(key string) bool {
	return m.ReturnValue
}

func (m *MockEnvGetter) GetIntEnv(key string) int {
	return m.ReturnInt
}

func (m *MockEnvGetter) GetDoubleEnv(key string) float64 {
	return m.ReturnFloat
}

func (m *MockEnvGetter) GetStrEnv(key string) string {
	return m.ReturnStr
}
