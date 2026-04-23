package cleaner

type Recovery struct {
	qm *QuarantineManager
}

func NewRecovery(qm *QuarantineManager) *Recovery {
	return &Recovery{qm: qm}
}

func (r *Recovery) RestoreFile(quarantinePath, originalPath string) error {
	return r.qm.Restore(quarantinePath, originalPath)
}
