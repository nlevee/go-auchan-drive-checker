# go-auchan-drive-checker

Ce script a pour objet la recherche de créneau dans les Drives Auchan.

## Usage

Pour l'identifiant de votre drive, on peut le trouver
dans les liens présent sur cette page : [Liste des drives](https://www.auchandrive.fr/drive/nos-drives/)

Le script va tourné en continue et va affiché sur la console si un créneau est disponible

Pour lancer le scripts :

```bash
./go-auchan-drive-checker -id [ID DU DRIVE]
```

Pour rendre accessible la recherche de créneau via une mini API :

```bash
./go-auchan-drive-checker -id [ID DU DRIVE] -port 8089 -host 0.0.0.0
```

Pour tester le serveur :

```bash
curl 127.0.0.1:8089/
```
