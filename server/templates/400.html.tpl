{% extends "layouts/not_login.html.tpl" %}

{% block title %}{{ title }}{% endblock %}

{% block content %}
<div class="error">
  <div class="main-board">
    <article>
      <div class="content">
        <h2>400 Bad Request.</h2>
        <h3>The request URL is invalid</h3>
      </div>
    </article>
  </div>
</div>
{% endblock %}
