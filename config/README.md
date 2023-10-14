config
======================

Тут лежит модуль, который загружает конфигурацию из окружения процесса

https://12factor.net/ru/config

Но, загружать какие-нибудь массивы и сложные структуры из окружения процесса
может быть сложно, поэтому лучше использовать какие-нибудь готовые инструменты
для хранения конфигов типа [etcd](https://github.com/etcd-io/etcd) или [Hashicorp Vault](https://github.com/hashicorp/vault)

Я знаю про https://github.com/ilyakaznacheev/cleanenv
и https://github.com/spf13/viper (там есть https://github.com/spf13/viper/issues/339),
но я не хочу так заморачиваться.
