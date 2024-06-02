ifdef OS
   RM = del /Q
   FixPath = $(subst /,\,$1)
else
   ifeq ($(shell uname), Linux)
      RM = rm -f
      FixPath = $1
   endif
endif

build:
	go build .
run:
	go run .
install:
   go mod download
clean:
	$(RM) $(call FixPath,queue.db)
	$(RM) $(call FixPath,logs/queue.log)
	$(RM) $(call FixPath,logs/queue.err)