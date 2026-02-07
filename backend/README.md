# Backend

Backend on toteutettu Go-kielellä ja vastaa ajanvarausjärjestelmän
liiketoimintalogiikasta, tietokantayhteyksistä sekä audit-logituksesta.

API on toteutettu HTTP:n yli, mutta rajapinta on määritelty
Protocol Buffer -tiedostoilla (proto-first-approach).

Backend noudattaa tuotantotason periaatteita, kuten selkeää vastuunjakoa,
vahvaa tyypitystä ja audit-logitusta.

---

## Vastuut

Backend vastaa seuraavista kokonaisuuksista:

- Käyttäjien hallinta (business user)
- API-rajapinta frontendille (HTTP)
- Tietoturva (salasanan hashays)
- Audit-logitus
- Tietokantainteraktiot (PostgreSQL)

---

## Teknologiat

- Go
- net/http
- Protocol Buffers
- protoc-gen-go
- protoc-gen-gohttp
- PostgreSQL
- bcrypt
- UUID (github.com/google/uuid)

---

## API-arkkitehtuuri

Backendin API on määritelty Protocol Buffer -tiedostoilla.
Näiden pohjalta generoidaan:

- vahvasti tyypitetyt request- ja response-rakenteet
- HTTP-handlerien runko

Koodin generointi tehdään käyttäen:
- `protoc-gen-go`
- `protoc-gen-gohttp`

Tämä lähestymistapa mahdollistaa:
- yhden lähteen totuudelle (proto-tiedostot)
- vähemmän manuaalista boilerplate-koodia
- selkeästi määritellyn ja tyypitetyn API:n
- helpon rajapinnan kehittämisen ja ylläpidon

API kommunikoi HTTP:n yli käyttäen Go:n `net/http`-kirjastoa.

---
## Arkkitehtuuri

Backend noudattaa kerroksellista rakennetta:

### Service-kerros
- Sisältää liiketoimintalogiikan
- Ei ole tietoinen tietokannan toteutuksesta
- Vastaa käyttötapausten orkestroinnista

### Repository-kerros
- Vastaa tietokantatoiminnoista
- Toteutettu interface-pohjaisesti
- Mahdollistaa tallennusratkaisun vaihtamisen

### Model-kerros
- Sisältää domain-mallit
- Ei riippuvuuksia infrastruktuuriin

Tämä rakenne parantaa:
- testattavuutta
- luettavuutta
- ylläpidettävyyttä

---

## Audit-logitus

Audit-logitus on toteutettu erillisenä servicenä (`audit.Service`).
Audit routet rekisteröinti tapahtuu for -loopissa, koska ajatellaan tämän tulevaisuudessa laajenevan.

### Periaatteet

- Audit-logi ei koskaan kaada requestia
- Audit-tallennus on erotettu omaksi repository-rajapinnakseen
- Mahdollistaa useita toteutuksia (esim. tietokanta, tapahtumapohjainen ratkaisu)

---
## Huom
Backend on rakennettu näyte- ja oppimisprojektina, mutta noudattaa tuotantotason periaatteita ja arkkitehtuuriratkaisuja.