# KAFP

Command `kafp` reads standard input and writes its lines into kafka topic.

## Example

Start kafka. Create a topic. Read the events. Following writes events into kafka topic.

```zsh
kafp -topic example -servers localhost
Hello Kafka
```
