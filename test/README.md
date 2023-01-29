To initialis4 eour CA and generate certs, we need to pass various config files
to `cfssl`. We need seperate config files to generate our CA and server certs and
we need a config file containing general config info about our CA

Files used:

ca-csr.json -> `cfssl` will use this file to configure oour CA's certificate
    - CN: Common Name, name for the CA
    - C: locality or municipality/city
    - ST: state or province
    - O: organisation
    - OU: organisational unit, such as department responsible for owning the key

ca-config.json -> use this file to define the CA's policy

server-csr.json -> `cfssl` will use it to configure our server's certificate
    the `hosts` field is a list of the domain names for which the certificate should be
    valid for. Since we're running our service locally we need "localost" and "127.0.0.1"