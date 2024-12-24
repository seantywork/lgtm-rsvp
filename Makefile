all:
	go build -o rsvp.out . 

vendor:

	apt-get update 

	apt-get install -y sqlite3

clean:
	rm -rf *.out 
