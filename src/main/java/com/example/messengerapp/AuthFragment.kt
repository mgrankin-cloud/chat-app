package com.example.messengerapp

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Button
import android.widget.EditText
import android.widget.TextView
import androidx.fragment.app.Fragment
import androidx.fragment.app.FragmentTransaction
import com.example.messengerapp.R


class AuthFragment : Fragment() {

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        val view = inflater.inflate(R.layout.fragment_auth, container, false)



        val authLine: TextView = view.findViewById(R.id.authLine)

        authLine.setOnClickListener {
            parentFragmentManager.beginTransaction()
                .replace(R.id.container, EmailInputFragment.newInstance())
                .addToBackStack(null) // Добавляем в back stack, чтобы можно было вернуться назад
                .setTransition(FragmentTransaction.TRANSIT_FRAGMENT_FADE) // Анимация перехода
                .commit()
        }

        val authButton: Button = view.findViewById(R.id.authButton)



        authButton.setOnClickListener{
            parentFragmentManager.beginTransaction()
                .replace(R.id.container, CodeVerificationFragment.newInstance(CodeVerificationFragment.FROM_LOGIN))
                .addToBackStack(null)
                .setTransition(FragmentTransaction.TRANSIT_FRAGMENT_FADE)
                .commit()

        }

        return view
    }

    companion object {
        fun newInstance(): AuthFragment {
            return AuthFragment()
        }
    }
}
