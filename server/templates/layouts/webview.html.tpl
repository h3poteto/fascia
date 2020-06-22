<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>{% block title %}Fascia{% endblock %}</title>
    <link rel="icon" href="/lp/images/favicon.ico" type="image/vnd.microsoft.icon">
    <link href='https://fonts.googleapis.com/css?family=Raleway' rel='stylesheet' type='text/css'>
    <link href="https://cdnjs.cloudflare.com/ajax/libs/octicons/4.4.0/font/octicons.css" rel="stylesheet" type="text/css">
    <link rel="stylesheet" href="/lp/css/lp-webview.css" media="all">
  </head>
  <body>
    <div id="content">
      {% block content %}{% endblock %}
    </div>
    <script src="https://use.fontawesome.com/080be9d465.js"></script>
    <script src="https://www.google.com/recaptcha/api.js?render=6Lf4lKcZAAAAAIHL6kGXMvMhmEFAJQvThnppcbZ9"></script>
    <script>
     grecaptcha.ready(function () {
       grecaptcha.execute('6Lf4lKcZAAAAAIHL6kGXMvMhmEFAJQvThnppcbZ9', { action: 'contact' }).then(function (token) {
         var recaptchaResponse = document.getElementById('recaptchaResponse');
         recaptchaResponse.value = token;
       });
     });
    </script>
  </body>
</html>
