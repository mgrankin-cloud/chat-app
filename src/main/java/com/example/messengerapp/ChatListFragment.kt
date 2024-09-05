package com.example.messengerapp

import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.Fragment
import androidx.recyclerview.widget.LinearLayoutManager
import androidx.recyclerview.widget.RecyclerView
import com.google.android.material.bottomnavigation.BottomNavigationView

class ChatListFragment : Fragment() {

    private lateinit var chatRecyclerView: RecyclerView
    private lateinit var bottomNavigationView: BottomNavigationView

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        val view = inflater.inflate(R.layout.activity_chatlist, container, false)

        // Инициализация RecyclerView
        chatRecyclerView = view.findViewById(R.id.chatRecyclerView)
        chatRecyclerView.layoutManager = LinearLayoutManager(context)
        chatRecyclerView.adapter = ChatListAdapter(getChatList())

        // Инициализация BottomNavigationView
        bottomNavigationView = view.findViewById(R.id.bottomNavigationView)
        bottomNavigationView.setOnNavigationItemSelectedListener { item ->
            when (item.itemId) {
                R.id.navigation_chats -> {
                    // Ничего не делаем, уже находимся на экране чатов
                    true
                }
                R.id.navigation_profile -> {
                    // Переход на экран профиля
                    parentFragmentManager.beginTransaction()
                        .replace(R.id.container, ProfileSettingsFragment.newInstance())
                        .addToBackStack(null)
                        .commit()
                    true
                }
                R.id.navigation_settings -> {
                    // Переход на экран настроек
                    parentFragmentManager.beginTransaction()
                        .replace(R.id.container, SettingsFragment.newInstance())
                        .addToBackStack(null)
                        .commit()
                    true
                }
                else -> false
            }
        }

        return view
    }

    private fun getChatList(): List<ChatItem> {
        // Замените на получение данных из вашего источника
        return listOf(
            ChatItem("Чат 1", "Последнее сообщение", "12:00"),
            ChatItem("Чат 2", "Привет!", "12:30"),
            ChatItem("Чат 3", "Как дела?", "13:00")
        )
    }
}
