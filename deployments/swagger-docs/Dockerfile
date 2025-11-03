FROM swaggerapi/swagger-ui:v5.10.5

COPY swagger /usr/share/nginx/html/specs

ENV URLS='[{"url":"/specs/auth/auth.swagger.json","name":"Auth Service"},{"url":"/specs/users/users.swagger.json","name":"Users Service"},{"url":"/specs/chat/chat.swagger.json","name":"Chat Service"},{"url":"/specs/social/social.swagger.json","name":"Social Service"},{"url":"/specs/notifications/notifications.swagger.json","name":"Notifications Service"}]'
ENV URLS_PRIMARY_NAME="Auth Service"

EXPOSE 8080

