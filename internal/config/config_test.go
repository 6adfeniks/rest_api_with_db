package config

import "testing"

func TestNewConfig(t *testing.T) {
	t.Run("right path", func(t *testing.T) {
		cfg, err := NewConfig("../../configs/config2.yml")
		if err != nil {
			t.Fatal(err)
		}

		if cfg.Database.User == "" ||
			cfg.Database.Password == "" ||
			cfg.Database.Dbname == "" ||
			cfg.Server.Port == "" || cfg.Server.Host == "" {
			t.Errorf("configs doesn't parses settings")
		}
	})

	t.Run("wrong path", func(t *testing.T) {
		_, err := NewConfig("www")
		if err.Error() != "open www: The system cannot find the file specified." {
			t.Errorf("The error got: %v, want: open www: The system cannot find the file specified.", err.Error())
		}
	})

	t.Run("wrong yml file", func(t *testing.T){
		_, err := NewConfig("../../testdata/fakeconfig.yml")
		if err.Error() != "yaml: line 8: could not find expected ':'"{
			t.Errorf("The error got: %v, want: yaml: line 8: could not find expected ':'", err.Error())
		}
	})


}
