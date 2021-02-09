FROM openjdk:8-jdk-alpine

ARG TZ="Asia/Seoul"

ENV TZ ${TZ}
ENV VERSION 1.0.0-SNAPSHOT
ENV APM_HOME /home/apm

RUN set -ex \
	&& chown ${RUN_USER}:${RUN_GROUP} ${APM_HOME} \
	&& ln -sf /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone \
	# Add apm user and disable sudo password checking
	&& useradd \
	  	--groups=sudo \
	  	--create-home \
	  	--home-dir=/home/apm \
	  	--shell=/bin/bash \
	  	--password=$(mkpasswd apm) \
	  	apm \
	&& sed -i '/%sudo[[:space:]]/ s/ALL[[:space:]]*$/NOPASSWD:ALL/' /etc/sudoers

# Setup environment
USER apm
WORKDIR ${APM_HOME}

COPY --from=builder apm-server/target/apm-server-${VERSION}.jar ${APM_HOME}/app.jar
COPY --from=builder apm-server/target/metatron.sh ${APM_HOME}

EXPOSE 8083

ENTRYPOINT ["java","-Dspring.profiles.active=prod","-jar","/app.jar"]

CMD ["./bin/metatron.sh"]