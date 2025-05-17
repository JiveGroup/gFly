{% extends "master.tpl" %}
    {% block body %}
    <p style="font-family: Helvetica, sans-serif; font-size: 16px; font-weight: normal; margin: 0; margin-bottom: 16px;">
        Hi {{ user_name }}
    </p>
    <p style="font-family: Helvetica, sans-serif; font-size: 16px; font-weight: normal; margin: 0; margin-bottom: 16px;">
        We change you greeting message. Please check below content
    </p>
    <p style="font-family: Helvetica, sans-serif; font-size: 16px; font-weight: normal; margin: 0; margin-bottom: 16px;">
        <b>{{ text }}</b>
    </p>
    <p style="font-family: Helvetica, sans-serif; font-size: 16px; font-weight: normal; margin: 0; margin-bottom: 16px;">
        Good luck! Hope it works.
    </p>
    {% endblock %}
