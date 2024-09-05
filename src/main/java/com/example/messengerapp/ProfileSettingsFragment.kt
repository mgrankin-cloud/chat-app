package com.example.messengerapp

import android.os.Bundle
import android.text.Editable
import android.text.TextWatcher
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.ImageView
import androidx.fragment.app.Fragment
import androidx.fragment.app.FragmentTransaction
import com.google.android.material.button.MaterialButton
import com.google.android.material.textfield.TextInputEditText

class ProfileSettingsFragment : Fragment() {

    private lateinit var nicknameEditText: TextInputEditText
    private lateinit var phoneEditText: TextInputEditText
    private lateinit var profileImageView: ImageView
    private lateinit var saveButton: MaterialButton

    private var originalNickname: String = "ТекущийНик"
    private var originalPhone: String = "+7 999 123-45-67"
    private var isPhotoChanged: Boolean = false

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        val view = inflater.inflate(R.layout.fragment_profile_settings, container, false)





        return view
    }


    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        nicknameEditText = view.findViewById(R.id.nicknameEditText)
        phoneEditText = view.findViewById(R.id.phoneEditText)
        profileImageView = view.findViewById(R.id.profileImage)
        saveButton = view.findViewById(R.id.saveButton)

        // Предзаполнение данными пользователя
        nicknameEditText.setText(originalNickname)
        phoneEditText.setText(originalPhone)

        // Изначально скрываем кнопку сохранения
        saveButton.visibility = View.GONE

        // Устанавливаем слушатели на изменения в текстовых полях
        nicknameEditText.addTextChangedListener(textWatcher)
        phoneEditText.addTextChangedListener(textWatcher)

        // Слушатель на изменение фото профиля
        profileImageView.setOnClickListener {
            // Логика изменения фото профиля
            isPhotoChanged = true
            checkForChanges()
        }
    }

    // Слушатель изменений текста
    private val textWatcher = object : TextWatcher {
        override fun beforeTextChanged(s: CharSequence?, start: Int, count: Int, after: Int) {}
        override fun onTextChanged(s: CharSequence?, start: Int, before: Int, count: Int) {
            checkForChanges()
        }
        override fun afterTextChanged(s: Editable?) {}
    }

    // Метод для проверки, были ли изменения
    private fun checkForChanges() {
        val isNicknameChanged = nicknameEditText.text.toString() != originalNickname
        val isPhoneChanged = phoneEditText.text.toString() != originalPhone

        if (isNicknameChanged || isPhoneChanged || isPhotoChanged) {
            saveButton.visibility = View.VISIBLE
        } else {
            saveButton.visibility = View.GONE
        }
    }
    companion object {
        fun newInstance(): ProfileSettingsFragment {
            return ProfileSettingsFragment()
        }
    }
}
