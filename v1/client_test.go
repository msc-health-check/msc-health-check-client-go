package v1

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestClientImpl_AddService(t *testing.T) {

	tests := []struct {
		Name                  string
		ProjectCheckRequest   ProjectCheckRequest
		ExpectError           error
		ExpectExistsID        bool
		ExpectAppName         string
		ExpectChecksOutSize   int
		ExpectLiveSignalsSize int
		ExpectErrorsSize      int
	}{
		{
			Name: "adicionando ProjectCheck com sucesso",
			ProjectCheckRequest: ProjectCheckRequest{
				URL: "https://google.com?q=\"https://https://msc-health-check.github.io/\"",
				AppName: uuid.NewString(),
			},
			ExpectError:           nil,
			ExpectExistsID:        true,
			ExpectChecksOutSize:   0,
			ExpectLiveSignalsSize: 1,
			ExpectErrorsSize:      0,
		},
		{
			Name:                  "adicionando ProjectCheck com erro: sem protocolo http e sem nome app",
			ProjectCheckRequest:   ProjectCheckRequest{},
			ExpectError:           errors.New("{\"error\":\"informe protocolo (http ou https),informe nome do app\"}"),
			ExpectAppName:         "",
			ExpectExistsID:        false,
			ExpectChecksOutSize:   0,
			ExpectLiveSignalsSize: 0,
			ExpectErrorsSize:      0,
		},
		{
			Name:                  "adicionando ProjectCheck com erro: sem nome app",
			ProjectCheckRequest:   ProjectCheckRequest{
				URL: "https://google.com?q=\"https://https://msc-health-check.github.io/\"",
			},
			ExpectError:           errors.New("{\"error\":\"informe nome do app\"}"),
			ExpectAppName:         "",
			ExpectExistsID:        false,
			ExpectChecksOutSize:   0,
			ExpectLiveSignalsSize: 0,
			ExpectErrorsSize:      0,
		},
		{
			Name:                  "adicionando ProjectCheck com erro: sem protocolo http",
			ProjectCheckRequest:   ProjectCheckRequest{
				URL: "google.com?q=\"https://https://msc-health-check.github.io/\"",
				AppName: uuid.NewString(),
			},
			ExpectError:           errors.New("{\"error\":\"informe protocolo (http ou https)\"}"),
			ExpectAppName:         "",
			ExpectExistsID:        false,
			ExpectChecksOutSize:   0,
			ExpectLiveSignalsSize: 0,
			ExpectErrorsSize:      0,
		},
	}

	for _, test := range tests {

		t.Run(test.Name, func(t *testing.T) {

			if len(test.ProjectCheckRequest.AppName) != 0 {
				test.ExpectAppName = test.ProjectCheckRequest.AppName
			}

			client := http.Client{
				Timeout: 20 * time.Minute,
			}
			clientImpl := NewClient(&client)
			projectCheck, err := clientImpl.AddService(test.ProjectCheckRequest)

			if err != nil {
				assert.Equal(t, test.ExpectError.Error(), err.Error())
			}

			if err == nil {
				assert.Equal(t, test.ExpectAppName, projectCheck.AppName)
				assert.Equal(t, test.ExpectErrorsSize, len(projectCheck.Errors))
				assert.Equal(t, test.ExpectChecksOutSize, len(projectCheck.ChecksOut))
				assert.Equal(t, test.ExpectExistsID, len(projectCheck.ID) > 0)
				assert.Equal(t, test.ExpectLiveSignalsSize, len(projectCheck.LiveSignals))
			}

		})
	}
}
