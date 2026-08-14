-- Project-local: lint Go against linux instead of the darwin host.
vim.lsp.config("gopls", { settings = { gopls = { env = { GOOS = "linux" } } } })
