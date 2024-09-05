package com.example.messengerapp

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.TextView
import androidx.fragment.app.Fragment
import androidx.fragment.app.FragmentTransaction
import com.google.android.material.button.MaterialButton
import com.google.android.material.textfield.TextInputEditText

class PasswordFragment : Fragment() {

    private lateinit var passwordEditText: TextInputEditText
    private lateinit var confirmPasswordEditText: TextInputEditText
    private lateinit var registerButton: MaterialButton

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        val view = inflater.inflate(R.layout.fragment_password, container, false)

        passwordEditText = view.findViewById(R.id.userPassword)
        confirmPasswordEditText = view.findViewById(R.id.userPasswordRepeat)
        registerButton = view.findViewById(R.id.regButton)






        registerButton.setOnClickListener {
            parentFragmentManager.beginTransaction()
                .replace(R.id.container, ProfileSettingsFragment.newInstance())
                .setTransition(FragmentTransaction.TRANSIT_FRAGMENT_FADE) // Анимация перехода
                .commit()
        }

        return view
    }

    companion object {
        fun newInstance(): PasswordFragment {
            return PasswordFragment()
        }
    }
}
