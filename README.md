## Prueba técnica en Go - AntPack

API REST para la gestión de usuarios con GO.

## Requerimientos

**GO 1.16**

**PostgreSQL 12, 13**

## Instalación en local

1. Instalar la versión de GO 1.16 o posterior de acuerdo a su sistema operativo.

   [Descargar GO](https://golang.org/dl/)


2. Clonar el proyecto.


3. Crear la base de datos.


4. Migraciones usando `golang-migrate` **(opcional)**.

    1. Descargar el ejecutable de `golang-migrate` para ejecutar las migraciones:

       [Descargar golang-migrate](https://github.com/golang-migrate/migrate/releases)

    2. En Windows guardar el archivo en la siguiente ruta:

       ```bash
       c:/go-migrate/migrate.exe
       ```
       Si elige otra ruta para guardar el ejecutable debe configurar el archivo `migrate.sh`

    3. Para ejecutar las migraciones abrir la consola `Git Bash` y usar el siguiente comando:

       ```bash
       export DBPWD=mypassword && . migrate.sh up
       ```
       Donde el valor de la variable de entorno `DBPWD` es la contraseña de la base de datos.

    4. Dentro de la carpeta [db](./migrations/db) de `migrations` encontrará el script para insertar registros de
       prueba.


5. Migraciones usando **pgAdmin**:

    1. Dentro de la carpeta [versions](./migrations/versions) de `migrations` encontrará el script para crear las tablas
       de la base de datos y en la carpeta [db](./migrations/db) el script para insertar registros de prueba.


6. Configurar variables de entorno en el archivo `.env` (Ver el archivo `.env.example`).


7. Ejecutar el servidor local:

   ```bash
    go run main.go
   ```

## Esquema

Ver [Migraciones](./migrations)

## Endpoints

La API REST proporciona las siguientes rutas:

### Auth

| Name | Route | Protected | Method |
| ---- | ----- | --------- | ------ |
| Signup | /signup | No | POST |
| Login | /login | No | POST |

### Users

| Name | Route | Protected | Method |
| ---- | ----- | --------- | ------ |
| Get all users  | /users | Yes | GET |
| Get User  | /users/:id | Yes | GET |

## Estructura del proyecto

Se utilizan los principios del patrón de diseño MVC para hacer mantenible y escalable el proyecto.

```
anpack-go
├── app                         Aplicación principal del proyecto.
│   ├── auth                    Servicio de autenticación.
│   ├── controllers             Funciones de controlador para una ruta en particular.
│   ├── middlewares             Middlewares que se utilizarán en el proyecto.
│   ├── models                  Tablas de la base de datos que se utilizarán como estructura de modelos.
│   ├── repositories            Persistir los datos en la base de datos.
│   ├── resources               Contiene estructuras que se usan como respuestas.
│   ├── ├── api                 Recursos asosiados a la API.
│   ├── ├── ├── errors          Tipos y manejo de errores.
│   ├── ├── ├── requests        Estructuras y otros modelos usados para hacer solicitudes.
│   ├── ├── ├── responses       Estructuras y otros modelos usadas como respuesta.
│   └── routes                  Rutas de la API.
│   └── services                Logica de negocio.
│   └── utils                   Funciones de ayuda utilizadas en todo el proyecto.
├── migrations                  Migraciones de base de datos
│   ├── db                      Archivos de consultas SQL.
│   ├── versions                Versiones individuales del esquema de la base de datos.
├── pkg                         Paquetes generales independientes del framework usado.
│   ├── config                  Archivos para leer variables de entorno.
│   ├── database                Archivos para conectarse a base de datos.
├── main.go                     Punto de entrada para iniciar el servidor.
├── .env.example                Archivo de ejemplo para configurar las variables de entorno.
├── migrate.sh                  Script para ejecutar las migraciones de la base de datos.
```

## Autor

**Luis Guillermo Gómez**

- [Github](https://github.com/luisgomez29)
