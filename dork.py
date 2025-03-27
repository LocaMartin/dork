#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from __future__ import print_function
import sys
import time

try:
    from googlesearch import search
except ImportError:
    print("\033[91m[ERROR] Missing dependency: googlesearch-python\033[0m")
    print("\033[93m[INFO] Install it using: pip install googlesearch-python\033[0m")
    sys.exit(1)

if sys.version_info[0] < 3:
    print("\n\033[91m[ERROR] This script requires Python 3.x\033[0m\n")
    sys.exit(1)

class Colors:
    RED = "\033[91m"
    BLUE = "\033[94m"
    GREEN = "\033[92m"
    YELLOW = "\033[93m"
    RESET = "\033[0m"

log_file = "dorks_output.txt"

def logger(data):
    with open(log_file, "a", encoding="utf-8") as file:
        file.write(data + "\n")

def dorks():
    global log_file
    try:
        dork = input(f"{Colors.BLUE}\n[+] Enter The Dork Search Query: {Colors.RESET}")
        
        user_choice = input(f"{Colors.BLUE}[+] Enter Total Number of Results (or 'all'): {Colors.RESET}").strip().lower()
        
        if user_choice == "all":
            total_results = 100  # Max results (Google's limit)
        else:
            try:
                total_results = int(user_choice)
                if total_results <= 0:
                    raise ValueError
            except ValueError:
                print(f"{Colors.RED}[ERROR] Enter a positive integer or 'all'{Colors.RESET}")
                return
        
        save_output = input(f"{Colors.BLUE}\n[+] Save Output? (Y/N): {Colors.RESET}").strip().lower()
        if save_output == "y":
            log_file = input(f"{Colors.BLUE}[+] Output Filename: {Colors.RESET}").strip()
            if not log_file:
                log_file = "dorks_output.txt"
            if not log_file.endswith(".txt"):
                log_file += ".txt"
        
        print(f"\n{Colors.GREEN}[INFO] Searching...{Colors.RESET}\n")
        
        fetched = 0
        try:
            # Use num_results instead of num, and remove the start parameter
            for result in search(dork, num_results=total_results):
                print(f"{Colors.YELLOW}[+] {Colors.RESET}{result}")
                if save_output == "y":
                    logger(result)
                fetched += 1
        except Exception as e:
            print(f"{Colors.RED}[ERROR] {str(e)}{Colors.RESET}")
        
        print(f"{Colors.GREEN}\n[✔] Found {fetched} results.{Colors.RESET}")
        
    except KeyboardInterrupt:
        print(f"\n{Colors.RED}[!] Exiting...{Colors.RESET}\n")
        sys.exit(1)

if __name__ == "__main__":
    dorks()
