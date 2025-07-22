package main

import (
	"fmt"
)

func commandMetatype(cfg *config, args ...string) error {
	metaType := "hash"
	if len(args) > 1 {
		return fmt.Errorf("usage : metatype [metatype]")
	} else if len(args) == 1 {
		metaType = args[0]
	}

	metaTypeResp, err := cfg.esapiClient.MetaTypes(metaType)
	if err != nil {
		return fmt.Errorf("error  %v :command_metatype.go in 18 lines", err)
	}

	// TODO: make text file with metatype respone or cache it
	fmt.Println(metaTypeResp)

	return nil
}
