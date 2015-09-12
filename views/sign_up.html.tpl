{% extends "layouts/login.html.tpl" %}

{% block title %}{{ title }}{% endblock %}

{% block content %}
<form action="/sign_up" method="post" role="form" name="sign_up" id="sign_up" class="pure-form pure-form-stacked">
    <fieldset>
        <label for="email">Email:</label>
        <input class="form-control" name="email" type="email" />

        <label for="password">Password:</label>
        <input class="form-control" name="password" type="password" />

        <label for="password-confirm">Password:</label>
        <input class="form-control" name="password-confirm" type="password" />

        <button class="pure-button pure-button-primary" type="submit">SignUp</button>
    </fieldset>
</form>
<p><a href="sign_in">SignIn</a></p>
{% endblock %}
