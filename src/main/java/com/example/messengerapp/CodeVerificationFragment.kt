package com.example.messengerapp

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import androidx.fragment.app.FragmentTransaction
import com.example.messengerapp.R
import com.google.android.material.button.MaterialButton
import com.google.android.material.textfield.TextInputEditText

class CodeVerificationFragment : Fragment() {

    private lateinit var codeEditTexts: Array<TextInputEditText>
    private lateinit var verifyButton: MaterialButton

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        val view = inflater.inflate(R.layout.fragment_code_verification, container, false)

        codeEditTexts = arrayOf(
            view.findViewById(R.id.codeInput1),
            view.findViewById(R.id.codeInput2),
            view.findViewById(R.id.codeInput3),
            view.findViewById(R.id.codeInput4)
        )

        verifyButton = view.findViewById(R.id.verifyButton)

        // Получение аргумента, переданного в фрагмент
        val source = arguments?.getString(ARG_SOURCE)

        // Установите слушатель нажатия
        verifyButton.setOnClickListener {
            val nextFragment = when (source) {
                FROM_REGISTRATION -> LoginPhoneFragment.newInstance()
                FROM_LOGIN ->  ProfileSettingsFragment.newInstance() // Переход на страницу профиля
                else -> LoginPhoneFragment.newInstance() // По умолчанию переход на страницу ввода телефона и логина
            }

            parentFragmentManager.beginTransaction()
                .replace(R.id.container, nextFragment)

                .setTransition(FragmentTransaction.TRANSIT_FRAGMENT_FADE)
                .commit()
        }

        return view
    }

    companion object {
        private const val ARG_SOURCE = "source"
        const val FROM_REGISTRATION = "fromRegistration"
        const val FROM_LOGIN = "fromLogin"

        fun newInstance(source: String): CodeVerificationFragment {
            val fragment = CodeVerificationFragment()
            val args = Bundle()
            args.putString(ARG_SOURCE, source)
            fragment.arguments = args
            return fragment
        }
    }
}
