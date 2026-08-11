package versions

func ValidateVersionNumber(versionNumber int) error {
	if versionNumber < 1 {
		return ErrInvalidVersion
	}

	return nil
}
