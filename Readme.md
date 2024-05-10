# Manual Técnico

  
A continuación se hace detalle de todos los componentes del proyecto:

### 1. Arquitectura del Sistema
En esta sección, se presenta la arquitectura general del sistema, incluyendo los componentes principales y cómo interactúan entre sí por medio de la red de Docker.
##### 1.1. Backend
API Rest realizada con GO, encargada de manejar el intercambio de datos entre el frontend y la base de datos. Así mismo, ejecuta los diversos endpoints que acceden a la información de los módulos de kernel. Esta tiene el puerto 8080 tanto interna como externamente.
##### 1.2. Frontend
Frontend realizado con React, apoyado de diversas librerías para la representación estadística del proyecto como gráficas de pie, entre otras. Entre las librerías utilizadas se tienen Bootstrap, Chart JS, Graphviz-React, y SweetAlert. Dicho componente expone el puerto 80 en la red de Docker, esto por medio de una configuración con NGINX. 
##### 1.3. Base de Datos
Se utilizó una imagen de Docker para instanciar una base de datos mySQL. Esta para almacenar los datos del uso histórico de la memoria RAM y del CPU, así como del historial de cada uno de los procesos creados por medio de la plataforma. Esta expone el puerto 3307 externamente y el puerto 3306 internamente. 
##### 1.4. Simulador de Estados
Esta es una página que hace uso de comandos Kill para comunicarse con el sistema operativo e iniciar, detener, continuar y terminar procesos. Así mismo, por medio de Graphviz se puede observar una representación de los estados que va tomando el proceso, así como el historial de estados por los cuales ha pasado dicho proceso.

##### 1.5. Módulos de Kernel
El módulo de kernel de la RAM se encarga de monitorear el uso de memoria RAM del sistema en tiempo real. Utiliza técnicas de acceso al hardware para recopilar datos sobre la cantidad de memoria utilizada y disponible, así como información detallada sobre los procesos que están utilizando la memoria en un momento dado. Esta información se utiliza para generar estadísticas y gráficos que muestran el uso de la memoria RAM a lo largo del tiempo.

Y el módulo de CPU es similar al módulo de monitoreo de RAM, este módulo se encarga de monitorear el uso de la CPU del sistema en tiempo real. Utiliza técnicas de acceso al hardware para recopilar datos sobre la carga de trabajo de la CPU, la utilización de núcleos individuales y la distribución de la carga entre los procesos en ejecución. Esta información se utiliza para generar estadísticas y gráficos que muestran el uso de la CPU a lo largo del tiempo.


### 2. Docker
La integración del proyecto con Docker permite encapsular y desplegar fácilmente cada uno de los componentes del sistema en contenedores independientes, proporcionando un entorno de desarrollo y ejecución consistente y reproducible. A continuación, se detallan los aspectos clave de la integración con Docker:

![conf](https://imgur.com/OcXruTW.png)

#### 2.1. Configuración del Docker Compose

El archivo `docker-compose.yml` define los servicios necesarios para el proyecto, incluyendo la base de datos MySQL, el backend en Go y el frontend en React. Cada servicio se configura con su imagen correspondiente, nombre del contenedor, puertos expuestos y dependencias entre servicios. Además, se utiliza el volumen `volume_mysql` para persistir los datos de la base de datos.

#### 2.2. Contenedor de Base de Datos MySQL

Se utiliza la imagen oficial de MySQL para instanciar un contenedor de base de datos, configurado con las variables de entorno necesarias para crear la base de datos y establecer la contraseña de root. El contenedor expone el puerto 3307 externamente para permitir la conexión desde el host y el puerto 3306 internamente para la comunicación entre contenedores.

#### 2.3. Contenedor de Backend en Go

Se utiliza una imagen personalizada del backend en Go, configurada para exponer el puerto 8080 y vinculada al contenedor de MySQL como dependencia. Además, se monta el directorio `/proc` del host en el contenedor para permitir el acceso a la información del sistema necesaria para el monitoreo de recursos.

#### 2.4. Contenedor de Frontend en React con NGINX

Se utiliza una imagen personalizada del frontend en React, configurada para exponer el puerto 80. Se utiliza NGINX como servidor web para servir los archivos estáticos del frontend y enrutar las solicitudes API al contenedor de backend. La configuración de NGINX se define en el archivo `nginx.conf`, que redirige las solicitudes `/api` al contenedor de backend.

#### 2.5. Interacción entre Contenedores

Los contenedores se comunican entre sí a través de la red interna de Docker, lo que permite que el frontend acceda al backend a través de la URL `http://mybackend:8080/api`. Esta configuración permite una comunicación eficiente y segura entre los diferentes componentes del sistema, facilitando el intercambio de datos y la colaboración entre servicios.

Con la integración con Docker, el proyecto se beneficia de la portabilidad, la escalabilidad y la consistencia del entorno de contenedores, lo que simplifica el desarrollo, el despliegue y la administración del sistema en diferentes entornos de ejecución.