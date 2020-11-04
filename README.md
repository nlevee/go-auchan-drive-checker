# go-auchan-drive-checker

* HS suite à la sortie du nouveau site *

Ce script a pour objet la recherche de créneau dans les Drives Auchan.

## Usage

Pour l'identifiant de votre drive, on peut le trouver dans les liens présent sur cette page : [Liste des drives](https://www.auchandrive.fr/drive/nos-drives/)

Si vous possédez une clé woosmap avec tous les drives vous pouvez renseigner ces informations dans un fichier .env (voir .env.dist)

Le script va tourner en continue et va afficher sur la console si un créneau est disponible

Pour lancer le scripts :

```bash
./go-auchan-drive-checker -id [ID DU DRIVE]
```

Ou si vous avez une clé woosmap, vous pouvez utiliser la recherche par code postal :

```bash
./go-auchan-drive-checker -cp [CODE POSTAL]
```

## Usage API

Pour rendre accessible la recherche de créneau via une mini API :

```bash
./go-auchan-drive-checker -port 8089 -host 0.0.0.0 &
```

Pour avoir la liste des clés de drive disponible :

```bash
curl 127.0.0.1:8089/stores?postalCode=[CODE POSTAL]
```

Pour ajouter un scrapper sur un store :

```bash
curl -XPUT 127.0.0.1:8089/scrappers/[ID DU DRIVE]
```

Pour checké l'état d'un drive :

```bash
curl 127.0.0.1:8089/scrappers/[ID DU DRIVE]
```
