package dotfiles

// func collectErrors(errs []error) error {
// 	if len(errs) == 0 {
// 		return nil
// 	}

// 	var errMsg string

// 	for _, err := range errs {
// 		errMsg += err.Error() + "\n"
// 	}

// 	return errors.New(errMsg)
// }

// func (p Profile) Init() error {
// 	for _, location := range p.Locations {
// 		err := os.Mkdir(location, os.ModePerm)
// 		if err != nil {
// 			return err
// 		}

// 		backend := loadBackend(p.Backend)
// 		err = backend.NewProfile(location)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
