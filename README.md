<p align="center"><a href="https://github.com/darkweak/dob"><img src="doc/logo.svg?sanitize=true" alt="dob logo"></a></p>

# Dob Table of Contents
1. [Dob multiple SSL files agregator](#project-description)
2. [Configuration](#configuration)

# dob multiple SSL certificates agregator

## Project description
Dob allow you to dynamically watch SSL certificates and keys files and get notified when one of them get an update.

## Configuration
dob is using a `config.yml` file for his own configuration. This file is listen in the code to avoid restart binary instance.  
You just have to create and map this file to dob instance to feel free to update it later.  
Then, the file should respect the format

```yaml
certificates: #Define certificates array
  - first: #Your certificate name
    domain: domain.com #Root domain linked to certificate
    SANs: # SANs are optional
      - api.domain.com #First subdomain
      - admin.domain.com #Second subdomain
    certificate: /path/certificate1 #Path to the certificate
    key: /path/key1 #Path to the key
  - first: #Your certificate name
    domain: domain2.com #Root domain linked to certificate
    SANs: # SANs are optional
      - *.domain.com #Every subdomains
    certificate: /path/certificate2 #Path to the certificate
    key: /path/key2 #Path to the key
  - third:
    domain: domain3.com
    certificate: /path/certificate3
    key: /path/key3
email: my@email.com #Fake email
output: test.json #Name of the output file name
```
