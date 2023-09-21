FROM confluentinc/cp-kafka:latest

USER root

COPY create-topic.sh /opt/kafka/create-topic.sh
RUN chmod +x /opt/kafka/create-topic.sh

CMD ["/opt/kafka/create-topic.sh"]