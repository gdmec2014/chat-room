### 说明

    还在开发中，请期待~

# chat-room  server

    chat-room server

# go version

    >= 1.12
    
# init
 
  go mod init    

# dev

    bee run -downdoc=true -gendoc=true
    
>若出现卡住了

    bee generate docs //可查看错误

# build

    go build

### 工具包

    cd  %GOPATH%\src\golang

    git clone https://github.com/golang/tools.git tools


    复制到 %GOPATH%\src\golang.org\x\tools

    go install github.com/ramya-rao-a/go-outline

    go install github.com/acroca/go-symbols

    go install golang.org/x/tools/cmd/guru

    go install golang.org/x/tools/cmd/gorename

    go install github.com/josharian/impl

    go install github.com/rogpeppe/godef

    go install github.com/sqs/goreturns

    go install github.com/golang/lint/golint

    go install github.com/cweill/gotests/gotests    

    单独处理golint,golint的源码位于https://github.com/golang/lint
    进入%GOPATH%\src\golang.org\x后执行git clone https://github.com/golang/lint下载golint需要的源码
    进入到%GOPATH%下，执行go install golang.org\x\lint

