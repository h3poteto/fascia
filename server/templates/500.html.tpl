{% extends "layouts/not_login.html.tpl" %}

{% block title %}{{ title }}{% endblock %}

{% block content %}
<div class="error">
  <div class="main-board">
    <article>
      <div class="content">
        <h2>Internal Server Error.</h2>
        <h3>We're sorry, but something went wrong.</h3>
      </div>
    </article>
  </div>
</div>
{% endblock %}
