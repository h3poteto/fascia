{% extends "layouts/not_login.html.tpl" %}

{% block title %}{{ title }}{% endblock %}

{% block content %}
<form action="/sign_in" method="post" role="form" name="sign_in" id="sign_in" class="pure-form pure-form-stacked">
    <fieldset>
        <label for="email">Email:</label>
        <input class="form-control" name="email" type="email" />

        <label for="password">Password:</label>
        <input class="form-control" name="password" type="password" />

        <input name="token" type="hidden" value="{{ token }}" />

        <button class="pure-button pure-button-primary" type="submit">SignIn</button>
    </fieldset>
</form>
<a href={{ oauthURL }}>SignIn with github</a>
<p><a href="sign_up">SignUp</a></p>
{% endblock %}
