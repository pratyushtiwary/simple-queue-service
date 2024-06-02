ifdef OS
   RM = del /Q
   FixPath = $(subst /,\,$1)
else
   ifeq ($(shell uname), Linux)
      RM = rm -f
      FixPath = $1
   endif
endif

install:
	go mod download
build:
	go build -mod=vendor .
run:
	go run .
clean:
	$(RM) $(call FixPath,queue.db)
	$(RM) $(call FixPath,logs/queue.log)
	$(RM) $(call FixPath,logs/queue.err)
