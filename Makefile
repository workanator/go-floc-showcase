all:
	@echo "Please use on of the targets from the list below."
	@echo "Examples:"
	@echo "  - 'ex_words' (Words)"
	@echo "  - 'ex_counter' (Counter)"

ex_words:
	@go run -race examples/words/main.go

ex_counter:
	@go run -race examples/counter/main.go
