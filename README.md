# Ajanvarausjärjestelmä
Ajanvarausjärjestelmä on full stack -projekti, joka kattaa frontendin, backendin, tietokannan sekä tietoturvan. Projekti demonstroi tuotantotason rakennetta, audit-logitusta ja backend–frontend-integraatiota.

Järjestelmä on suunniteltu kahdelle käyttäjäryhmälle:
- **Business user** (yritysasiakkaat)
- **Customer** (kuluttaja-asiakkaat)

## Käyttäjäroolit

**Business userit** ovat yritysasiakkaita, jotka esittäytyvät palvelussa ja
avaavat kalentereistaan varattavia aikoja. He voivat toimia yksinyrittäjinä
(esim. kampaaja, hieroja) tai osana yritystä (yrityksen omistaja tai työntekijä), jossa on useita työntekijöitä.

**Customer-käyttäjät** ovat kuluttaja-asiakkaita, jotka varaavat aikoja
yritysasiakkailta.


## Teknologiat

### Backend
- Go
- Proto Buffer
- (go)HTTP
- PostgresSQL
- Audit-logitus

### Frontend
- Next.js
- TypeScript
- API-integraatio

## Ominaisuudet
- Käyttäjän (business user) luonti gRPC:n kautta
- PostgreSQL-pohjainen repository
- Audit-logitus (ei kaada requestia)
- Salasanan turvallinen hashays (bcrypt)
- Selkeä service / repository -jako

## Miksi tämä projekti?

Projekti on rakennettu osoittamaan kykyäni suunnitella ja toteuttaa full stack -järjestelmää tuotantotason periaatteita noudattaen.

Erityistä huomiota on kiinnitetty:
- audit-logitukseen
- tietoturvaan
- virheiden hallintaan
- vastuiden selkeään jakamiseen koodissa

---

## Lisätietoa

Backendin tarkempi tekninen kuvaus löytyy
niiden omista README-tiedostoista. Frontendin kuvaus tulee myöhemmin, kunhan sitä aletaan kehittämään.