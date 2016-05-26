SOURCES := $(wildcard *.go)
BIN := alfred-github-jump
FILES := $(BIN) info.plist icon.png

Github\ Jump.alfredworkflow: $(FILES)
	zip -j "$@" $^

alfred-github-jump: $(SOURCES)
	go build -o alfred-github-jump $(SOURCES)

clean:
	rm $(BIN) Github\ Jump.alfredworkflow