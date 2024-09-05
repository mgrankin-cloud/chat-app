package com.example.messengerapp

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.TextView
import androidx.fragment.app.Fragment
import androidx.fragment.app.FragmentTransaction
import com.example.messengerapp.R
import com.google.android.material.button.MaterialButton
import com.google.android.material.textfield.TextInputEditText

class LoginPhoneFragment : Fragment() {

    private lateinit var loginEditText: TextInputEditText
    private lateinit var phoneEditText: TextInputEditText
    private lateinit var nextButton: MaterialButton

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        val view = inflater.inflate(R.layout.fragment_login_phone, container, false)

        loginEditText = view.findViewById(R.id.userLogin)
        phoneEditText = view.findViewById(R.id.userPhone)
        nextButton = view.findViewById(R.id.nextButton)



            // Установите слушатель нажатия
            nextButton.setOnClickListener {
                parentFragmentManager.beginTransaction()
                    .replace(R.id.container, PasswordFragment.newInstance())
                    .addToBackStack(null) // Добавляем в back stack, чтобы можно было вернуться назад
                    .setTransition(FragmentTransaction.TRANSIT_FRAGMENT_FADE) // Анимация перехода
                    .commit()
            }

            return view
        }

    companion object {
        fun newInstance(): LoginPhoneFragment {
            return LoginPhoneFragment()
        }
    }
}
