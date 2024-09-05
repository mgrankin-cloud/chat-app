package com.example.messengerapp

import android.os.Bundle
import androidx.appcompat.app.AppCompatActivity
import androidx.fragment.app.Fragment


class RegistrationActivity : AppCompatActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        // Фрагмент ввода почты
        supportFragmentManager.beginTransaction()
            .replace(R.id.fragment_container, EmailInputFragment())
            .commit()
    }

    fun navigateToCodeVerification(email: String) {
        // Переход к фрагменту подтверждения кода
        supportFragmentManager.beginTransaction()
            .replace(R.id.fragment_container, CodeVerificationFragment())
            .addToBackStack(null)
            .commit()
    }

    fun navigateToLoginAndPhone() {
        // Переход к фрагменту логина и телефона
        supportFragmentManager.beginTransaction()
            .replace(R.id.fragment_container, LoginPhoneFragment())
            .addToBackStack(null)
            .commit()
    }

    fun navigateToPassword() {
        // Переход к фрагменту пароля
        supportFragmentManager.beginTransaction()
            .replace(R.id.fragment_container, PasswordFragment())
            .addToBackStack(null)
            .commit()
    }

    fun completeRegistration(password: String, confirmPassword: String) {
        // Логика завершения регистрации
    }

    companion object {
        fun newInstance(): RegistrationActivity {
return RegistrationActivity()
        }
    }
}
