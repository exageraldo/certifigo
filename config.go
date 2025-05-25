package certifigo

import (
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

func NewEnvCredentials() (*EnvCredentials, error) {
	credentials := &EnvCredentials{}
	credentials.bindEnvs()

	viper.SetTypeByDefaultValue(true)
	viper.AutomaticEnv()

	err := viper.Unmarshal(&credentials)
	if err != nil {
		return nil, err
	}

	return credentials, nil
}

type EnvCredentials struct {
	EmailSender   string `mapstructure:"EMAIL_SENDER"`
	EmailPassword string `mapstructure:"EMAIL_PASSWORD"`

	JWTToken string `mapstructure:"JWT_TOKEN"`
}

func (c *EnvCredentials) CheckEmailCredentials() bool {
	if c.EmailSender == "" || c.EmailPassword == "" {
		return false
	}
	return true
}

func (c EnvCredentials) bindEnvs() {
	st := reflect.TypeOf(c)
	for t := 0; t < st.NumField(); t++ {
		tag := st.Field(t).Tag.Get("mapstructure")
		name := strings.Split(tag, ",")[0]
		if name == "" {
			// if a name is not passed, we will ignore it
			continue
		}
		// the only error that occurs is when no value is passed to the "BindEnv" function,
		// so we can ignore it.
		_ = viper.BindEnv(name)
	}
}
