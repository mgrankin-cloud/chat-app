    package com.example.messengerapp

    import android.os.Bundle
    import android.view.LayoutInflater
    import android.view.View
    import android.view.ViewGroup
    import android.widget.Button
    import android.widget.TextView
    import androidx.fragment.app.Fragment
    import androidx.fragment.app.FragmentTransaction
    import com.example.messengerapp.R
    import com.google.android.material.button.MaterialButton
    import com.google.android.material.textfield.TextInputEditText

    class EmailInputFragment : Fragment() {

            private lateinit var emailEditText: TextInputEditText
            private lateinit var nextButton: MaterialButton




        override fun onCreateView(
            inflater: LayoutInflater, container: ViewGroup?,
            savedInstanceState: Bundle?
        ): View? {
            val view = inflater.inflate(R.layout.fragment_email, container, false)

            // Найдите TextView с ID authLine
            val nextButton: TextView = view.findViewById(R.id.nextButton)

            // Установите слушатель нажатия
            nextButton.setOnClickListener {
                parentFragmentManager.beginTransaction()
                    .replace(R.id.container, CodeVerificationFragment.newInstance(CodeVerificationFragment.FROM_REGISTRATION))
                    .addToBackStack(null)
                    .setTransition(FragmentTransaction.TRANSIT_FRAGMENT_FADE)
                    .commit()
            }

            return view
        }


        companion object {
            fun newInstance(): EmailInputFragment {
            return EmailInputFragment()
            }
        }
    }
