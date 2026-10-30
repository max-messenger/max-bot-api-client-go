# Установка в системeу, обновление глобального списока доверенных центров сертификации.

## Скачиваем актуальные корневой и выпускающий сертификаты напрямую

Эти ссылки ведут на инфраструктуру Госуслуг/НУЦ

```sh
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates wget


mkdir /usr/local/share/ca-certificates/russian-trusted/

wget https://gu-st.ru/content/lending/russian_trusted_root_ca_pem.crt -O russian_trusted_root_ca.crt 
wget https://gu-st.ru/content/lending/russian_trusted_sub_ca_pem.crt -O russian_trusted_sub_ca.crt


update-ca-certificates

```
