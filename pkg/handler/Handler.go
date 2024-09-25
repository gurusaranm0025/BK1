package handler

type Handler struct {
	InputFiles   []string
	InputFolders []string
	OutputFiles  []string
}

func (h *Handler) Pack() {}

func (h *Handler) UnPack() {}
