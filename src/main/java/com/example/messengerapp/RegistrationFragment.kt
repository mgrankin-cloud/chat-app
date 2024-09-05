package com.example.messengerapp

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.TextView
import androidx.fragment.app.Fragment
import androidx.fragment.app.FragmentTransaction

class RegistrationFragment  : Fragment() {
    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        // Inflate the layout for this fragment
        val view = inflater.inflate(R.layout.fragment_reg, container, false)
        val authLine2: TextView = view.findViewById(R.id.authLine)

        // Установите слушатель нажатия
        authLine2.setOnClickListener {
            parentFragmentManager.beginTransaction()
                .replace(R.id.container, AuthFragment.newInstance())
                .addToBackStack(null) // Добавляем в back stack, чтобы можно было вернуться назад
                .setTransition(FragmentTransaction.TRANSIT_FRAGMENT_FADE) // Анимация перехода
                .commit()
        }

        return view
    }

    companion object {
        fun newInstance(): RegistrationFragment {
            return RegistrationFragment()
        }
    }
}