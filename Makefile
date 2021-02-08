OUTPUT=cmd/aphrodite-web.runtime

GITSHA=`git rev-parse --short HEAD`
GITTAG=`git describe --abbrev=0 --tags`
BUILDTIME=`date +%FT%T%z`
LDFLAGS=-ldflags "-X main.GitSha=${GITSHA} -X main.GitTag=${GITTAG} -X main.BuildTime=${BUILDTIME} -s -w"

build: 
	GOOS=linux go build ${LDFLAGS} -o ${OUTPUT} startup.go