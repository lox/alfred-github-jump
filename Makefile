SOURCES := $(wildcard *.go)
BIN := alfred-github-jump
FILES := $(BIN) info.plist icon.png

build: Github\ Jump.alfredworkflow

Github\ Jump.alfredworkflow: $(FILES)
	zip -j "$@" $^

alfred-github-jump: $(SOURCES)
	CGO_ENABLED=1 go build -o alfred-github-jump $(SOURCES)

clean:
	-rm $(BIN) Github\ Jump.alfredworkflow
