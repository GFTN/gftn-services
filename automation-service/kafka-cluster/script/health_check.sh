
CA_CERT=${CA_CERT:-ibmca.crt}
CA_KEY=${CA_KEY:-ibmca.key}
CA_PASSWORD=${CA_PASSWORD:-WorldWire-test}
STORE_PASS=${STORE_PASS:-Worldwire-teststore}
STORE_LOCATION="/var/private/ssl"
LIST=kafak-cli
OU=${OU:-IBMBlockchain}
O=${O:-IBM}
L=${L:-SG}
C=${C:-SG}
ENVIRONMENT=${ENVIRONMENT:-test}
DOMAIN=${DOMAIN:-ww}
DIR="$(cd "$(dirname "$0")" && pwd)"


#pull ca cert, key and password
SECRET_NAME=/$ENVIRONMENT/$DOMAIN/kafka_ca_cert
source $DIR/secret_pull.sh | grep SecretString |awk '{print $2}' | tr -d '",' > temp
python -m base64 -d  temp > ./$CA_CERT

SECRET_NAME=/$ENVIRONMENT/$DOMAIN/kafka_ca_key
source $DIR/secret_pull.sh | grep SecretString |awk '{print $2}' | tr -d '",' > temp
python -m base64 -d  temp >  ./$CA_KEY

SECRET_NAME=/$ENVIRONMENT/$DOMAIN/kafka_ca_password
source $DIR/secret_pull.sh | grep SecretString |awk '{print $2}' | tr -d '",' > temp
python -m base64 -d  temp >  passtemp
CA_PASSWORD=$(cat passtemp)

source $DIR/generate_sign_store_kafka.sh


source $DIR/generate_sign_store
echo "security.protocol=SSL
ssl.truststore.location=/var/private/ssl/kafka-cli/kafka.kafka-cli.truststore.jks 
ssl.truststore.password=$STORE_PASS
ssl.keystore.location=/var/private/ssl/kafka-cli/kafka.kafka-cli.keystore.jks 
ssl.keystore.password=$STORE_PASS
ssl.key.password=$STORE_PASS" > /temp
seq 42 | kafka-console-producer --broker-list kafka-1:19092 --topic test --producer.config /temp
# todo
kafka-console-consumer --bootstrap-server kafka_2:9092 --from-beginning --max-messages 42 --topic test --consumer.config /temp

echo "health check pass"