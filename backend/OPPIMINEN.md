Tämä tiedosto kokoaa keskeiset opit ajanvarausjärjestelmäprojektista. Se toimii muistilistana ja kertauksena.

---

## 1. Handler → Service → Repository

**Handler:**
- Vastaa requestin vastaanottamisesta ja parametrien välittämisestä service-kerrokselle.
- Aseta default-arvot (esim. `limit = 50`, `offset = 0`).
- Ei pidä tehdä liiketoimintalogiikkaa tai tietokantakyselyjä.

**Service:**
- Käsittelee liiketoimintalogiikan.
- Muuntaa domain-mallit protoiksi.
- Ei tiedä HTTP/RPC-transporttia.

**Repository:**
- Ainoa kerros, joka kommunikoi tietokannan kanssa.
- Hoitaa lopulliset tarkistukset ja rajoitukset (esim. max limit, offset >= 0).
- Toteuttaa paginationin (LIMIT/OFFSET) ja tarvittaessa filttereitä.

---

## 2. Pagination ja proto-mallit
- Käytä **Request/Response** protoja, älä suoraa listaa (esim. UserList).
  
```proto
message GetUsersRequest {
    int32 limit = 1;
    int32 offset = 2;
}

message GetUsersResponse {
    repeated User users = 1;
}
```

- Handler lukee requestista limit/offset ja asettaa default-arvot.
- Repository tekee lopullisen rajauksen SQL:ssä.
- Tämä rakenne mahdollistaa myös filttereiden ja hakuehtojen lisäämisen tulevaisuudessa.

---
## 3. ActorID ja audit-logit

- Käytä context.Contextia tallentamaan actorID (kuka teki toimenpiteen).
- Middleware (auth) asettaa actorID:n contextiin.
- Audit-service hakee actorID:n suoraan contextista:
```
if id, ok := actorctx.ActorIDFromContext(ctx); ok {
    actorID = id
}
```
- Tämä eriyttää audit-logituksen infrastruktuurista ja tekee koodista testattavampaa.
---
## 4. Refaktorointi ja commit käytännöt
- Pidä commitit pieninä ja fokusoituina yhteen asiaan.
- Refaktoroi koodi ennen uusien ominaisuuksien lisäämistä.
- Käytä selkeitä commit-viestejä, jotka kuvaavat muutoksen tarkoituksen.

---
## 5. Kerrosten vastuiden erottelu

- Handler: RPC/HTTP <→ Service
- Service: Liiketoimintalogiikka, domain-mallien käsittely
- Repository: Tietokanta, SQL-queryt, data-muunnokset, validointi

Muista: Älä sekoita kerroksia. Tämä helpottaa testattavuutta ja ylläpidettävyyttä.

---
## 6. Hyvät käytännöt
- Default-arvot requestissa, lopullinen validointi repositoryssa.
- Käytä Request/Response protoja yhtenäisen rajapinnan varmistamiseksi.
- Vältä string-keytä contextissa, käytä omaa type key struct{}.
- Eri kerrokset helpottavat unit-testien kirjoittamista.
- Lisää kontekstia virheilmoituksiin, älä paljasta tietokannan yksityiskohtia frontille.
- Rajaa SQL-queryissä limit, offset ja mahdolliset hakuehdot.
- Handler pysyy minimalistisena, service hoitaa business- ja audit-logiikan, repository hoitaa DB:n.